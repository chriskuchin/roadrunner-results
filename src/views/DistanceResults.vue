<template>
  <div class="container mt-5">
    <div class="mb-5">
      <div class="field has-addons">
        <p class="control">
          <a :class="['button', 'is-static', 'is-large']">
            <!-- attempts/allowed_attempts -->
            3/3
          </a>
        </p>
        <p class="control is-expanded has-icons-left">
          <input :class="['input', 'is-large']" ref="input" :type="getInputType" placeholder="Bib Number" />
          <span class="icon is-large is-left">
            <icon :icon="['fas', 'fa-user']"></icon>
          </span>
        </p>
        <p class="control" v-if="false">
          <a :class="['button', 'is-static', 'is-large']">
            <!-- best attempt -->
            <!-- {{ time }} -->
          </a>
        </p>
      </div>
      <label class="checkbox">
        <input type="checkbox" v-model="letterInput">
        Allow Letters
      </label>
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
        <div class="control is-expanded has-icons-left">
          <input class="input is-large" placeholder="Feet" type="number">
          <span class="icon is-large is-left">
            <icon :icon="['fas', 'fa-ruler']"></icon>
          </span>
        </div>
        <div class="control is-expanded has-icons-left">
          <input class="input is-large" placeholder="Inches" type="number">
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
      <div class="level-item has-text-centered" v-for="attempt in attempts">
        <div>
          <p class="heading">
            <icon :icon="['fas', 'fa-crown']"></icon> Attempt #{{ attempt.attempt }}
          </p>
          <p class="title">{{ formatCentimeters(attempt.distance, format) }}</p>
        </div>
      </div>
    </div>
    <!-- <p :class="['help']" v-if="result.show">{{ result.msg }}</p> -->
    <fab button_type="is-danger">
      <icon icon="fa-solid fa-stopwatch"></icon>
    </fab>
  </div>
</template>

<script>
import { setAuthHeader } from '../api/auth';
import { formatCentimeters } from '../utilities'
import FAB from '../components/Fab.vue'

export default {
  components: {
    'fab': FAB,
  },
  mounted: function () {
  },
  data: function () {
    return {
      format: "ftin",
      letterInput: false,
      attempts: [
        {
          attempt: 3,
          distance: 162.56,
        },
        {
          attempt: 1,
          distance: 157.48,
        },
        {
          attempt: 2,
          distance: 154.94,
        }
      ]
    }
  },
  computed: {
    getInputType() {
      return this.letterInput ? "text" : "number"
    }
  },
  methods: {
    recordAttempt: function () {

    },
    formatCentimeters,
  }
};
</script>
