<template>
  <div class="container">
    <div class="section mt-1 mb-0 py-0" v-if="this.$route.name != 'events'">
      <router-link :to="'/races/' + this.$route.params.raceId + '/events'">
        <icon icon="fa-solid fa-arrow-left-long"></icon> Back to Events
      </router-link>
    </div>
    <router-view></router-view>
  </div>
</template>

<script>
import Breadcrumb from "../components/Breadcrumb.vue";
import { useRaceStore } from "../store/race";

const store = useRaceStore()
export default {
  components: {
    breadcrumb: Breadcrumb,
  },
  data: function () {
    return {
    }
  },
  mounted: function () {
    store.loadRace(this.$route.params.raceId)
  },
  methods: {},
  computed: {
    backPath: function () {
      return this.$route.path.slice(0, this.$route.path.lastIndexOf('/'))
    },
    paths: function () {
      var breadcrumb = [{
        path: "/races",
        display: "Races"
      }]
      var path = "/races"
      var pathSegments = this.$route.path.replace("/races/", "").split("/")
      for (let i = 0; i < pathSegments.length; i++) {
        var display = pathSegments[i]
        path += "/" + pathSegments[i]
        if (display == "" || display == "events")
          continue

        if (i == 0 && store.getName != "") {
          display = store.getName
        } else if (i == 2) {
          display = store.eventName(pathSegments[i])
        }

        breadcrumb.push({
          path: path,
          display: display
        })
      }

      return breadcrumb;
    },
  },
};
</script>