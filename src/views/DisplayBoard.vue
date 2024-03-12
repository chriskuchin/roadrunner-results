<template>
  <div class="section">
    <h1 class="title">{{ title }}</h1>
    <div class="table-container">
      <table class="table mx-auto is-size-2" v-if="page.results.length > 0">
        <thead>
          <tr>
            <th>Place</th>
            <th>Time</th>
            <th>Bib</th>
            <th>First Name</th>
            <th>Last Name</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="result, place in page.results">
            <th>{{ place + 1 }}</th>
            <td>{{ formatResult(result.result_ms) }}</td>
            <td>{{ result.bib_number }}</td>
            <td>{{ result.first_name }}</td>
            <td>{{ result.last_name }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
import { getEventResults } from '../api/events';
import { listTimers } from '../api/timers';
import { formatMilliseconds } from '../utilities';


export default {
  mounted: async function () {
    await this.refreshPages()
    this.currentPage = this.pages.shift()

    const that = this
    this.timeout = setInterval(() => {
      that.currentPage = that.pages.shift()
      if (that.pages.length == 0) {
        that.refreshPages()
      }
    }, 30000)
  },
  data: function () {
    return {
      pages: [],
      currentPage: null,
      timeout: null,
      page: {
        start: 0,
        results: [],
      }
    }
  },
  watch: {
    async currentPage() {
      if (this.currentPage == null)
        return

      this.page.results = []
      const raceId = this.$route.params.raceId
      const eventId = this.$route.params.eventId
      switch (this.currentPage.type) {
        case "timer":
          this.page.results = await getEventResults(raceId, eventId, "", "", "", "", [this.currentPage.id], "desc")
          break;

        default:
          var gender = ""
          for (let filter of this.currentPage.filters) {
            if (filter.key == "gender") {
              gender = filter.value
            }
          }
          this.page.results = await getEventResults(raceId, eventId, [], [gender], [], [], [], "desc")
          break;
      }
    }
  },
  computed: {
    title() {
      if (this.currentPage != null && this.currentPage.title != "") {
        return this.currentPage.title
      }

      return "Race Info"
    }
  },
  methods: {
    formatResult: function (result) {
      return formatMilliseconds(result)
    },
    refreshPages: async function () {
      const allTimers = await listTimers(this.$route.params.raceId, this.$route.params.eventId)

      this.pages = []
      let heat = 1
      for (let timer of allTimers) {
        if (timer.timer_start > 0)
          this.pages.push({
            type: "timer",
            id: timer.id,
            title: `Heat ${heat}`
          })

        heat++
      }

      this.pages.push({
        type: "result",
        title: "Results Female",
        filters: [{
          key: "gender",
          value: "Female"
        }]
      })

      this.pages.push({
        type: "result",
        title: "Results Male",
        filters: [{
          key: "gender",
          value: "Male"
        }]
      })
    }
  }
};
</script>