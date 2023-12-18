import { defineStore } from "pinia";

export const useErrorBus = defineStore("error", {
	state: () => ({
		show: false,
		msg: "",
	}),
	getters: {
		errorMsg: (state) => state.msg,
		showNotification: (state) => state.show,
	},
	actions: {
		hide: function () {
			this.show = false;
		},
		handle: function (msg) {
			this.msg = msg;
			this.show = true;
		},
	},
});
