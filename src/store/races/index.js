import { defineStore } from "pinia";
import { getRaces } from "../../api/races";
import { ref } from 'vue'

export const useRacesStore = defineStore('races', () => {
	const races = ref([])

	function loadRaces() {
		getRaces().then((result) => {
			races.value = result
		})
	}

	return { races, loadRaces }
})
