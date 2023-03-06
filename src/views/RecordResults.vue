<template>
  <div>
    <div class="section">
      <div class="tabs">
        <ul>
          <li @click="tabSelect('manual')" :class="{ 'is-active': isActiveTab('manual') }">
            <a>Manual</a>
          </li>
          <li @click="tabSelect('scanner')" :class="{ 'is-active': isActiveTab('scanner') }">
            <a>Scanner</a>
          </li>
        </ul>
      </div>
      <div class="container">
        <result-scanner v-if="isActiveTab('scanner')" />
        <result-input v-else-if="isActiveTab('manual')" />
      </div>
      <results-table />
    </div>
  </div>
</template>

<script>
import Scanner from "../components/Scanner.vue";
import RacerInput from "../components/ResultsInput.vue";
import ResultsTable from "../components/ResultsTable.vue";
import Breadcrumb from "../components/Breadcrumb.vue";

export default {
  components: {
    "result-scanner": Scanner,
    "result-input": RacerInput,
    "results-table": ResultsTable,
    breadcrumb: Breadcrumb,
  },
  data: function () {
    return {
      activeTab: "manual",
    };
  },
  methods: {
    isActiveTab: function (tab) {
      return this.activeTab == tab;
    },
    tabSelect: function (tab) {
      this.activeTab = tab;
    },
  },
  computed: {
    paths: function () {
      return [
        {
          path: "/races",
          display: "Races",
        },
        {
          path: "/" + this.$route.params.raceId,
          display: "Race 123",
        },
        {
          path: "/" + this.$route.params.raceId + "/results",
          display: "Results",
        },
      ];
    },
  },
};
</script>