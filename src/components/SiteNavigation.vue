<template>
  <nav class="navbar is-dark" role="navigation" aria-label="main navigation">
    <div class="navbar-brand">
      <a class="navbar-item" href="/">
        <img v-bind:src="logo_url" height="30px" alt="RoadRunner Results" title="results.roadrunners.club" />
      </a>
      <a role="button" class="navbar-burger" aria-label="menu" aria-expanded="false" data-target="navbarBasicExample"
        @click="toggleBurger" :class="{ 'is-active': activeBurger }">
        <span aria-hidden="true"></span>
        <span aria-hidden="true"></span>
        <span aria-hidden="true"></span>
      </a>
    </div>
    <div id="navbarBasicExample" class="navbar-menu" :class="{ 'is-active': activeBurger }">
      <div class="navbar-start">
        <!-- <router-link to="/" class="navbar-item">Races</router-link>
        <router-link to="/" class="navbar-item">Results</router-link>
        <div class="navbar-item has-dropdown is-hoverable">
          <a class="navbar-link"> Manage </a>
          <div class="navbar-dropdown">
            <router-link to="/" class="navbar-item">Register</router-link>
            <router-link to="/" class="navbar-item">Record</router-link>
            <router-link to="/" class="navbar-item">Timer</router-link>
          </div>
        </div> -->
      </div>
      <div class="navbar-end">
        <div class="navbar-item" v-if="!isLoggedIn">
          <div class="buttons">
            <router-link to="/signup" class="button is-primary">
              <strong>Sign Up</strong>
            </router-link>
            <router-link to="/login" class="button is-light">
              Log in
            </router-link>
          </div>
        </div>
        <div class="navbar-item has-dropdown is-hoverable" v-else>
          <a class="navbar-link"> {{ userDisplayName }} </a>
          <div class="navbar-dropdown">
            <router-link to="/" class="navbar-item">Profile</router-link>
            <router-link to="/" class="navbar-item">Settings</router-link>
            <hr class="dropdown-divider" />
            <a @click="logout" class="navbar-item">Sign Out</a>
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>

<script>
import logo from "../assets/images/logo.png";
import { mapState, mapActions } from "pinia";
import { useUserStore } from "../store/user";

export default {
  components: {
  },
  data: function () {
    return {
      activeBurger: false,
      apiKeyModal: false,
      apiKey: ""
    };
  },
  methods: {
    toggleApiKeyModal: function () {
      this.apiKeyModal = !this.apiKeyModal
    },
    toggleBurger: function () {
      this.activeBurger = !this.activeBurger;
    },
    saveToken: function () {
      saveAPIToken(this.apiKey)
      this.toggleApiKeyModal()
    },
    ...mapActions(useUserStore, ['logout'])
  },
  computed: {
    ...mapState(useUserStore, ['isLoggedIn', 'userDisplayName']),
    logo_url: function () {
      return logo;
    },
  },
};
</script>
