//go:build with_pulse

#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <pulse/pulseaudio.h>

#include "volume_pulse.h"

typedef struct volume_pulse {
	pa_mainloop *mainloop;
	pa_context *context;
	char *default_sink;
	uint32_t sink_index;
	pa_volume_t prev_avg;
	int prev_muted;
	char *muted_label;
	uintptr_t handle;
} volume_pulse_t;

static void server_info_cb(pa_context *, const pa_server_info *, void *);
static void sink_info_cb(pa_context *, const pa_sink_info *, int, void *);

extern void update(const char *, uintptr_t);

static void
context_state_cb(pa_context *c, void *userdata)
{
	volume_pulse_t *v = userdata;
	static pa_subscription_mask_t m = PA_SUBSCRIPTION_MASK_SERVER |
	                                  PA_SUBSCRIPTION_MASK_SINK;

	switch (pa_context_get_state(c)) {
	case PA_CONTEXT_READY:
		pa_operation *op;
		if ((op = pa_context_subscribe(c, m, NULL, NULL)))
			pa_operation_unref(op);
		if ((op = pa_context_get_server_info(v->context, server_info_cb, v)))
			pa_operation_unref(op);
		break;
	case PA_CONTEXT_FAILED:
		/* FALLTHROUGH */
	case PA_CONTEXT_TERMINATED:
		fprintf(stderr, "volume_pulse: connection lost\n");
		update("", v->handle);
		pa_mainloop_quit(v->mainloop, 1);
		break;
	default:
		/* NOP */
		break;
	}
}

static void
server_info_cb(pa_context *c, const pa_server_info *i, void *userdata)
{
	volume_pulse_t *v = userdata;
	if (!i || !i->default_sink_name) {
		fprintf(stderr, "volume_pulse: server_info_cb: no default sink\n");
		return;
	}

	fprintf(stderr, "volume_pulse: default sink is now %s\n", i->default_sink_name);

	free(v->default_sink);
	v->default_sink = strdup(i->default_sink_name);
	if (!v->default_sink)
		return;

	pa_operation *op = pa_context_get_sink_info_by_name(c, v->default_sink,
	                                                    sink_info_cb, v);
	if (op)
		pa_operation_unref(op);
}

static void
sink_info_cb(pa_context *c, const pa_sink_info *i, int eol, void *userdata)
{
	(void)c;
	if (eol != 0 || !i)
		return;

	volume_pulse_t *v = userdata;
	v->sink_index = i->index;

	pa_volume_t avg = i->volume.channels > 0
		? pa_cvolume_avg(&i->volume)
		: PA_VOLUME_INVALID;

	if (v->prev_muted != -1 && v->prev_muted == (int)i->mute &&
	    v->prev_avg != PA_VOLUME_INVALID && avg != PA_VOLUME_INVALID &&
	    v->prev_avg == avg)
		return;

	v->prev_muted = i->mute;
	v->prev_avg = avg;

	if (i->mute) {
		update(v->muted_label, v->handle);
		return;
	}
	if (i->volume.channels == 0) {
		update("", v->handle);
		return;
	}

	/*
	max_perc =
		PA_VOLUME_MAX*100 / PA_VOLUME_NORM = 214748364700/65536 = 3276799
	*/
	unsigned perc = (unsigned)((uint64_t)avg * 100 / PA_VOLUME_NORM);

	char buf[9]; /* "3276799%\0" */
	snprintf(buf, sizeof(buf), "%u%%", perc);
	update(buf, v->handle);
}

static void
subscribe_cb(pa_context *c, pa_subscription_event_type_t t, uint32_t index,
             void *userdata)
{
	volume_pulse_t *v = userdata;

	pa_subscription_event_type_t type = t & PA_SUBSCRIPTION_EVENT_TYPE_MASK;
	if (type != PA_SUBSCRIPTION_EVENT_CHANGE)
		return;

	pa_subscription_event_type_t facility = t & PA_SUBSCRIPTION_EVENT_FACILITY_MASK;
	if (facility == PA_SUBSCRIPTION_EVENT_SERVER) {
		if (pa_context_get_state(v->context) == PA_CONTEXT_READY) {
			pa_operation *op = pa_context_get_server_info(v->context,
			                                              server_info_cb, v);
			if (op)
				pa_operation_unref(op);
		}
		return;
	}

	if (facility == PA_SUBSCRIPTION_EVENT_SINK &&
	    v->sink_index != 0 && v->sink_index == index) {
		pa_operation *op = pa_context_get_sink_info_by_index(c, index,
		                                                     sink_info_cb, v);
		if (op)
			pa_operation_unref(op);
	}
}

volume_pulse_t *
volume_pulse_new(const char *muted_label, uintptr_t h)
{
	volume_pulse_t *v = calloc(1, sizeof(*v));
	if (!v)
		return NULL;

	v->prev_avg    = PA_VOLUME_INVALID;
	v->prev_muted  = -1;
	v->muted_label = strdup(muted_label);
	v->handle      = h;

	if (!v->muted_label) {
		free(v);
		return NULL;
	}

	return v;
}

void
volume_pulse_run_mainloop(volume_pulse_t *v)
{
	useconds_t delay = INITIAL_RETRY_DELAY_US;

	for (;;) {
		v->mainloop = pa_mainloop_new();
		if (!v->mainloop)
			goto retry;

		v->context = pa_context_new(pa_mainloop_get_api(v->mainloop), CLIENT_NAME);
		if (!v->context) {
			pa_mainloop_free(v->mainloop);
			v->mainloop = NULL;
			goto retry;
		}

		pa_context_set_state_callback(v->context, context_state_cb, v);
		pa_context_set_subscribe_callback(v->context, subscribe_cb, v);

		if (pa_context_connect(v->context, NULL, 0, NULL) < 0) {
			fprintf(stderr, "volume_pulse: pa_context_connect: %s\n",
			        pa_strerror(pa_context_errno(v->context)));
			pa_context_unref(v->context);
			v->context = NULL;
			pa_mainloop_free(v->mainloop);
			v->mainloop = NULL;
			goto retry;
		}

		delay = INITIAL_RETRY_DELAY_US;

		int retval;
		pa_mainloop_run(v->mainloop, &retval);
		fprintf(stderr, "volume_pulse: mainloop exited, retval=%d\n", retval);

		pa_context_unref(v->context);
		pa_mainloop_free(v->mainloop);
		free(v->default_sink);

		v->mainloop     = NULL;
		v->context      = NULL;
		v->default_sink = NULL;
		v->sink_index   = 0;
		v->prev_avg     = PA_VOLUME_INVALID;
		v->prev_muted   = -1;

	retry:
		update("", v->handle);
		fprintf(stderr, "volume_pulse: retrying in %luus\n",
		        (unsigned long)delay);
		usleep(delay);
		delay = delay < MAX_RETRY_DELAY_US / 2
			? delay * 2
			: MAX_RETRY_DELAY_US;
	}
}
