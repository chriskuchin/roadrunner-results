import { defineStore } from "pinia";
import { setAuthHeader } from "../../api/auth";

export const useMediaStore = defineStore("media", {
	state: () => ({
		videoStream: null,
		videoElement: null,
		recordVideo: false,
		recorder: null,
		chunks: [],
		videoQuality: [
			{
				width: 1920,
				height: 1080,
			},
			{
				width: 1280,
				height: 720,
			},
			{
				width: 960,
				height: 540,
			},
		],
		settings: {
			audio: false,
			video: {
				width: 1920,
				height: 1080,
				focusMode: "continuous",
				frameRate: 60,
				facingMode: {
					ideal: "environment",
				},
			},
		},
	}),
	getters: {},
	actions: {
		async startCamera(video) {
			this.videoElement = video;
			this.videoStream = await navigator.mediaDevices.getUserMedia(
				this.settings,
			);
			video.srcObject = this.videoStream;
			video.play();
		},
		async recordVideo() {
			if (this.videoStream) {
				const options = {
					videoBitsPerSecond: 2500000000,
					mimeType: "video/webm;codecs=vp8",
				};
				this.recorder = new MediaRecorder(this.videoStream, options);
				this.recorder.ondataavailable = (e) => {
					this.chunks.push(e.data);
				};

				this.recorder.onstop = () => {
					const blob = new Blob(this.chunks, { type: "video/webm;codecs=vp8" });
					const download = document.createElement("a");
					download.href = URL.createObjectURL(blob);
					download.download = "raceVideo.webm";

					download.click();
					this.chunks = [];
					this.recorder = null;
				};
				this.recorder.start();
			}
		},
		async saveVideo() {
			if (this.recorder) {
				this.recorder.stop();
			}
		},
		async stopCamera(video) {
			const tracks = this.videoStream.getTracks();

			for (let i = 0; i < tracks.length; i++) {
				const track = tracks[i];
				track.stop();
			}

			video.srcObject = null;
			this.videoStream = null;
		},
		async takePicture(raceID, eventID, finishTime, elapsedTime) {
			if (this.videoStream) {
				const track = this.videoStream.getVideoTracks()[0];
				const canvas = document.createElement("canvas");
				const context = canvas.getContext("2d");

				canvas.width = track.getSettings().width;
				canvas.height = track.getSettings().height;

				context.drawImage(this.videoElement, 0, 0, canvas.width, canvas.height);
				canvas.toBlob((imageBlob) => {
					this.uploadFinishPhoto(
						imageBlob,
						raceID,
						eventID,
						finishTime,
						elapsedTime,
					);
				}, "image/png");
			}
		},
		async uploadFinishPhoto(imgBlob, raceID, eventID, finishTime, elapsedTime) {
			const url = `/api/v1/races/${raceID}/events/${eventID}/results`;
			const formData = new FormData();
			formData.append(
				"photo-finish",
				imgBlob,
				`finisher_${finishTime}_${elapsedTime}.png`,
			);

			const res = await fetch(
				url,
				await setAuthHeader({
					method: "POST",
					body: formData,
				}),
			);

			if (res.ok) {
				// need error handling added to this method
				console.log("Success!!!");
			}
		},
		async load() {
			navigator.mediaDevices.enumerateDevices().then((devices) => {
				// Better Devince handling
				console.log(devices);
			});
		},
	},
});
