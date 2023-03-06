<template>
  <div>
    <table class="table">
      <thead>
        <tr>
          <th><abbr>Race</abbr></th>
          <th><abbr>Participants</abbr></th>
          <th><abbr>Events</abbr></th>
          <th><abbr>Results</abbr></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="race in races" :key="race.id">
          <td><router-link :to="'/races/' + race.id + '/'">{{ race.name }}</router-link></td>
          <td><router-link :to="'/races/' + race.id + '/participants'">Participants</router-link></td>
          <td><router-link :to="'/races/' + race.id + '/events'">Events</router-link></td>
          <td><router-link :to="'/races/' + race.id + '/results'">Results</router-link></td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
export default {
  data: function () {
    return {
      races: [],
    };
  },
  mounted: function () {
    this.getRaces();
  },
  methods: {
    getRaces: function () {
      fetch("/api/v1/races", {
        method: "GET",
      })
        .then((res) => {
          return res.json();
        })
        .then((res) => {
          console.log(res)
          this.races = res;
        });
    },
  },
};
</script>