<template>
  <div class="container">
    <aside class="menu primary">
      <p class="menu-label">General</p>
      <ul class="menu-list">
        <li><a>Dashboard</a></li>
        <li><a>Customers</a></li>
      </ul>
      <p class="menu-label">Administration</p>
      <ul class="menu-list">
        <li><a>Team Settings</a></li>
        <li>
          <a class="is-active">Manage Your Team</a>
          <ul>
            <li><a>Members</a></li>
            <li><a>Plugins</a></li>
            <li><a>Add a member</a></li>
          </ul>
        </li>
        <li><a>Invitations</a></li>
        <li><a>Cloud Storage Environment Settings</a></li>
        <li><a>Authentication</a></li>
      </ul>
      <p class="menu-label">Transactions</p>
      <ul class="menu-list">
        <li><a>Payments</a></li>
        <li><a>Transfers</a></li>
        <li><a>Balance</a></li>
      </ul>
    </aside>

    <nav class="breadcrumb mt-3 is-centered mb-0" aria-label="breadcrumbs">
      <ul>
        <li :class="{ 'is-active': !isEventPage }">
          <router-link :to="'/races/' + this.$route.params.raceId + '/events'">
            {{ raceName }}
          </router-link>
        </li>
        <li :class="{ 'is-active': isEventPage }" v-if="isEventPage" class="is-active">
          <a href="#" aria-current="page">
            {{ eventName }}
          </a>
        </li>
      </ul>
    </nav>
    <router-view></router-view>
  </div>
</template>

<script>
import { useRaceStore } from "../store/race";
import { mapState, mapActions } from "pinia";

export default {
  components: {},
  data: function () {
    return {
    }
  },
  mounted: function () {
    this.loadRace()
    this.loadEvents()
  },
  methods: {
    ...mapActions(useRaceStore, ['loadRace', 'loadEvents']),
  },
  computed: {
    ...mapState(useRaceStore, {
      raceName: (store) => store.name,
      eventName: 'eventName'
    }),
    getEventId: function () {
      if (this.$route.params.eventId && this.$route.params.eventId !== "")
        return this.$route.params.eventId
      else if (this.$route.query.eventId && this.$route.query.eventId !== "")
        return this.$route.query.eventId
    },
    backPath: function () {
      return this.$route.path.slice(0, this.$route.path.lastIndexOf('/'))
    },
    isEventPage: function () {
      return (this.$route.params.eventId && this.$route.params.eventId !== "") || (this.$route.query.eventId && this.$route.query.eventId !== "")
    }
  },
};
</script>

<style scoped>
/* Styling for the side drawer */
.menu.primary {
  padding-top: 25px;
  position: fixed;
  top: 52px;
  left: 0;
  bottom: 0;
  width: 250px;
  z-index: 1000;
  /* Ensure the drawer appears on top of other content */
  background: var(--bulma-body-background-color);
  transform: translateX(-100%);
  transition: transform 0.3s ease-in-out;
}

.menu.primary.is-active {
  transform: translateX(0);
}

/* Styling for the button
.button {
  position: fixed;
  top: 20px;
  left: 20px;
} */
</style>