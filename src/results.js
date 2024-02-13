import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";

import { createPinia } from "pinia";

import "./css/results.scss";

/* import the fontawesome core */
import { library } from "@fortawesome/fontawesome-svg-core";

/* import font awesome icon component */
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";

import { onAuthStateChanged } from "firebase/auth";
import { auth } from "./firebase";

import { useUserStore } from "./store/user";

if ("serviceWorker" in navigator) {
	window.addEventListener("load", () => {
		navigator.serviceWorker
			.register("/service-worker.js")
			.then((registration) => {
				console.log("SW registered: ", registration);
			})
			.catch((registrationError) => {
				console.log("SW registration failed: ", registrationError);
			});
	});
}

/* import specific icons */
import {
	faArrowLeftLong,
	faDownload,
	faEllipsisV,
	faFileCsv,
	faFlagCheckered,
	faPlay,
	faPlus,
	faRepeat,
	faRuler,
	faStopwatch,
	faUser,
	faCrown,
	faAngleDown
} from "@fortawesome/free-solid-svg-icons";
library.add(
	faPlus,
	faStopwatch,
	faRuler,
	faRepeat,
	faPlay,
	faDownload,
	faFileCsv,
	faEllipsisV,
	faArrowLeftLong,
	faFlagCheckered,
	faUser,
	faCrown,
	faAngleDown,
);

require("./assets/");

let initialized = false;
const pinia = createPinia();
const app = createApp(App)
	.component("icon", FontAwesomeIcon)
	.use(router)
	.use(pinia);

onAuthStateChanged(auth, (user) => {
	if (user) {
		useUserStore().loadUser(user);
	}

	if (!initialized) {
		initialized = true;
		app.mount("#app");
	}
});
