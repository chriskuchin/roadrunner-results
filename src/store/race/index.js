import { defineStore } from 'pinia'
import { useRoute } from 'vue-router'
import { ref, computed } from 'vue'
import { getRace, getRaceVolunteers } from '../../api/races'
import { getParticipants } from '../../api/participants'
import { getRaceEvents } from '../../api/events'

export const useRaceStore = defineStore('race', () => {
  const route = useRoute()

  const id = ref(route.params.raceId)
  const name = ref('')
  const ownerId = ref('')
  const events = ref([])
  const participants = ref([])
  const volunteers = ref([])

  const getRaceId = computed(() => route.params.raceId)
  const getParticipant = computed(() => (bib) => participants.value.find((entry) => entry.bibNumber === bib.toString()))

  const participantFirstName = computed(() => (bib) => {
    const participant = participants.value.find((entry) => entry.bibNumber === bib.toString())
    if (participant)
      return participant.firstName

    return '-'
  })
  const participantLastName = computed(() => (bib) => {
    const participant = participants.value.find((entry) => entry.bibNumber === bib.toString())
    if (participant)
      return participant.lastName

    return '-'
  })
  const participantBirthYear = computed(() => (bib) => {
    const participant = participants.value.find((entry) => entry.bibNumber === bib.toString())
    if (participant)
      return participant.birthYear

    return '-'
  })

  const eventName = computed(() => {
    const event = events.value.find((entry) => {
      return entry.eventId === route.params.eventId
    })
    if (event) {
      return event.description
    }
    return "-"
  })

  function loadRace() {
    if (id !== route.params.raceId) {
      getRace(route.params.raceId).then((race) => {
        id.value = race.id
        name.value = race.name
        ownerId.value = race.owner_id
      })
    }
  }

  function loadParticipants() {
    getParticipants(route.params.raceId, 500, 0).then((result) => {
      participants.value = result
    })
  }

  function loadEvents() {
    getRaceEvents(route.params.raceId).then((result) => {
      events.value = result
    })
  }

  function loadVolunteers() {
    getRaceVolunteers(route.params.raceId).then((result) => {
      volunteers.value = result
    })
  }

  return { id, name, ownerId, events, participants, volunteers, loadRace, loadParticipants, loadEvents, loadVolunteers, getParticipant, eventName, getRaceId, participantFirstName, participantLastName, participantBirthYear }
})
