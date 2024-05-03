<template>
  <div class="block mt-5 mx-5">
    <div class="mb-3">
      <div class="field has-addons">
        <p class="control">
          <a :class="['button', 'is-static', 'is-large']">
            {{ attempts.length }}
          </a>
        </p>
        <p :class="['control', 'is-expanded', 'has-icons-left', processing ? 'is-loading' : '']">
          <input :class="['input', 'is-large']" ref="input" :type="getInputType" placeholder="Bib Number"
            v-on:blur="lookupAthlete" v-model="athlete.bib" />
          <span class="icon is-large is-left">
            <icon :icon="['fas', 'fa-user']"></icon>
          </span>
        </p>
        <p class="control" v-if="false">
          <a :class="['button', 'is-static', 'is-large']"></a>
        </p>
      </div>
      <label class="checkbox">
        <input type="checkbox" v-model="letterInput">
        Allow Letters
      </label>
    </div>
    <div class="level mb-5" v-if="athlete.info.firstName !== ''">
      <div class="level-item has-text-centered">
        <div>
          <p class="heading">{{ athlete.info.team }}</p>
          <h1 class="title">{{ athlete.info.firstName }} {{ athlete.info.lastName }}</h1>
          <p>({{ athlete.info.birthYear }} {{ athlete.info.gender }})</p>
        </div>
      </div>
    </div>
    <div class="mb-5">
      <div class="field has-addons">
        <!-- <div class="control">
          <div class="dropdown is-hoverable">
            <div class="dropdown-trigger">
              <button class="button is-large is-info" aria-haspopup="true" aria-controls="distance-measuring-menu">
                <span>Distance Measuring</span>
                <span class="icon is-small">
                  <icon :icon="['fas', 'fa-angle-down']" aria-hidden="true"></icon>
                </span>
              </button>
            </div>
            <div class="dropdown-menu" id="distance-measuring-menu" role="menu">
              <div class="dropdown-content">
                <a href="#" class="dropdown-item">
                  Feet & Inches
                </a>
                <a href="#" class="dropdown-item">
                  Meters & Centimeters
                </a>
              </div>
            </div>
          </div>
        </div> -->
        <div :class="['control', 'is-expanded', 'has-icons-left', processing ? 'is-loading' : '']">
          <input :class="['input', 'is-large']" placeholder="Feet" type="number" ref="l-msr">
          <span class="icon is-large is-left">
            <icon :icon="['fas', 'fa-ruler']"></icon>
          </span>
        </div>
        <div :class="['control', 'is-expanded', 'has-icons-left', processing ? 'is-loading' : '']">
          <input :class="['input', 'is-large']" placeholder="Inches" type="number" ref="s-msr">
          <span class="icon is-large is-left">
            <icon :icon="['fas', 'fa-ruler']"></icon>
          </span>
        </div>
      </div>
      <div class="control">
        <label class="radio">
          <input type="radio" name="format" value="ftin" v-model="format">
          Feet & Inches
        </label>
        <label class="radio">
          <input type="radio" name="format" value="cm" v-model="format">
          Centimeters
        </label>
        <label class="radio">
          <input type="radio" name="format" value="mcm" v-model="format">
          Meters & Centimentes
        </label>
      </div>
    </div>
    <div class="level">
      <div class="level-item has-text-centered" v-for="(attempt, key) in attempts">
        <div>
          <p class="heading">
            <icon v-if="key == 0" :icon="['fas', 'fa-crown']"></icon> Attempt #{{ attempt.attempt_number }}
          </p>
          <p class="title">{{ formatCentimeters(attempt.result, format) }}</p>
        </div>
      </div>
    </div>
    <!-- <p :class="['help']" v-if="result.show">{{ result.msg }}</p> -->
    <fab button_type="is-danger" @click="recordAttempt">
      <icon icon="fa-solid fa-stopwatch"></icon>
    </fab>
  </div>
</template>

<script>
import { formatCentimeters, calculateCentimeters } from '../utilities'
import { getParticipantByBib } from '../api/participants'
import { listEventAttempts, recordEventAttempt } from '../api/attempts'
import FAB from '../components/Fab.vue'
import { useErrorBus } from '../store/error'
import { mapActions } from 'pinia'

export default {
  components: {
    'fab': FAB,
  },
  unmounted: function () {
    window.removeEventListener('keypress', this.keyPressEvent)
  },
  mounted: function () {
    window.addEventListener('keypress', this.keyPressEvent)
  },
  data: function () {
    return {
      processing: false,
      athlete: {
        info: {
          firstName: "",
          lastName: "",
          team: "",
          gender: "",
          grade: "",
          birthYear: "",
        },
        bib: ""
      },
      current_attempt: 0,
      format: "ftin",
      letterInput: false,
      attempts: []
    }
  },
  computed: {
    getInputType() {
      return this.letterInput ? "text" : "number"
    }
  },
  methods: {
    ...mapActions(useErrorBus, { 'handleError': 'handle' }),
    keyPressEvent: function (e) {
      if (e.key == "Enter") {
        this.recordAttempt()
        this.$refs['l-msr'].focus()
      }
    },
    lookupAthlete: async function (e) {
      this.processing = true
      if (this.athlete.bib !== "") {
        let participant = await getParticipantByBib(this.$route.params.raceId, this.athlete.bib)

        if (Object.keys(participant).length > 0) {
          this.athlete.info.firstName = participant.first_name
          this.athlete.info.lastName = participant.last_name
          this.athlete.info.team = participant.team
          this.athlete.info.gender = participant.gender
          this.athlete.info.birthYear = participant.birth_year
        }

        this.attempts = await listEventAttempts(this.$route.params.raceId, this.$route.params.eventId, this.athlete.bib)
      } else {
        this.athlete.info.firstName = ""
        this.athlete.info.lastName = ""
        this.athlete.info.team = ""
        this.athlete.info.gender = ""
        this.athlete.info.birthYear = ""

        this.attempts = []
      }
      this.processing = false
    },
    recordAttempt: async function () {
      this.processing = true
      const largeVal = Number(this.$refs['l-msr'].value)
      const smallVal = Number(this.$refs['s-msr'].value)
      const distance = calculateCentimeters(largeVal, smallVal, this.format)

      recordEventAttempt(this.$route.params.raceId, this.$route.params.eventId, this.athlete.bib, this.attempts.length + 1, distance).then(() => {
        this.$refs['l-msr'].value = ""
        this.$refs['s-msr'].value = ""
        this.processing = false
      }).catch((err) => {
        this.processing = false
        this.handleError(err)
      })
    },
    formatCentimeters,
  }
};
</script>