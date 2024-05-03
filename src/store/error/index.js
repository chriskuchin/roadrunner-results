import { defineStore } from "pinia";

export const useErrorBus = defineStore("error", {
	state: () => ({
		msg: [],
	}),
	getters: {
		hasError: (state) => {
			return state.msg.length > 0
		},
		errorCount: (state) => state.msg.length,
		errorMsg: (state) => state.msg[0],
		showNotification: (state) => state.show,
	},
	actions: {
		hide: function () {
			this.msg.pop()

			if (this.msg.length === 0) {
				this.show = false;
			}
		},
		handle: function (msg) {
			this.msg.push(msg);
			this.show = true;
		},
	},
});
