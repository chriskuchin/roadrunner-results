<template>
  <div>
    <div class="level mb-4">
      <div class="level-item has-text-centered">
        <div>
          <p class="heading">Participants</p>
          <p class="title">{{ totalParticipants }}</p>
        </div>
      </div>
      <div class="level-item has-text-centered">
        <div>
          <p class="heading">Finishers</p>
          <p class="title">{{ totalFinishers }}</p>
        </div>
      </div>
      <div class="level-item has-text-centered">
        <div>
          <p class="heading">Finish Rate</p>
          <p class="title">{{ (totalFinishers / totalParticipants * 100).toFixed(2) }}%</p>
        </div>
      </div>
      <!-- <div class="level-item has-text-centered">
        <div>
          <p class="heading">Events</p>
          <p class="title">{{ events.length }}</p>
        </div>
      </div> -->
    </div>
    <div class="grid gap-3 is-col-min-3">
      <div class="cell">
        <div class="box">
          <p class="title">Participants</p>
          <canvas id="birth-year"></canvas>
        </div>
      </div>
      <div class="cell">
        <div class="box">
          <p class="title">Events</p>
          <ul class="has-text-centered">
            <li v-for="event in events" :key="event.eventId">
              <router-link :to="'/races/' + this.$route.params.raceId + '/divisions?eventsId=' + event.eventId">
                {{ event.description }} ({{ event.finishers }})
              </router-link>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Chart from 'chart.js/auto'
import { useRaceStore } from "../store/race";
import { mapState } from "pinia";
import { RouterLink } from 'vue-router';

const raceState = useRaceStore()

export default {
  components: {
    "router-link": RouterLink
  },
  mounted: function () {
    if (this.maleValues.length > 0) {
      this.updateCharts()
    }
    raceState.$subscribe(this.updateCharts)
  },
  unmounted: function () {
    this.birthYearChart = null
  },
  data: function () {
    return {
      event: {
        description: ""

      }
    }
  },
  methods: {
    updateCharts: function () {
      if (this.birthYearChart == null) {
        this.birthYearChart = new Chart(document.getElementById('birth-year'),
          {
            type: 'bar',
            data: {
              labels: this.yearLabels,
              datasets: [
                {
                  label: 'Male Athletes',
                  data: this.maleValues
                },
                {
                  label: 'Female Athletes',
                  data: this.femaleValues
                },
              ]
            },
            options: {
              scales: {
                x: {
                  stacked: true
                },
                y: {
                  stacked: true
                }
              }
            }
          }
        );
      }
    }
  },
  computed: {
    ...mapState(useRaceStore, ['yearLabels', 'maleValues', 'femaleValues', 'totalParticipants', 'totalFinishers', 'eventTotal', 'events'])
  }
};
</script>