<template>
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
  </div>
</template>

<script>
import Scanner from "../components/Scanner.vue";
import RacerInput from "../components/ResultsInput.vue";

export default {
  components: {
    "result-scanner": Scanner,
    "result-input": RacerInput,
  },
  data: function () {
    return {
      activeTab: "manual",
    };
  },
  created: function () { },
  methods: {
    isActiveTab: function (tab) {
      return this.activeTab == tab;
    },
    tabSelect: function (tab) {
      this.activeTab = tab;
    },
  },
  computed: {},
};
</script>