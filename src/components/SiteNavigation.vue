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
        <router-link to="/" class="navbar-item">Races</router-link>
        <router-link to="/" class="navbar-item">Results</router-link>
        <div class="navbar-item has-dropdown is-hoverable">
          <a class="navbar-link"> Manage </a>
          <div class="navbar-dropdown">
            <router-link to="/" class="navbar-item">Register</router-link>
            <router-link to="/" class="navbar-item">Record</router-link>
            <router-link to="/" class="navbar-item">Timer</router-link>
          </div>
        </div>
      </div>
      <div class="navbar-end">
        <div class="navbar-item">
          <div class="buttons">
            <a @click="toggleApiKeyModal" class="button is-link">
              <strong>API Key</strong>
            </a>
            <router-link to="/signup" class="button is-primary">
              <strong>Sign Up</strong>
            </router-link>
            <router-link to="/login" class="button is-light">
              Log in
            </router-link>
          </div>
        </div>
      </div>
    </div>
    <modal :show="apiKeyModal" @close="toggleApiKeyModal">
      <p class="title">Input API Key</p>
      <div class="field">
        <label class="label">API Key</label>
        <div class="control">
          <input class="input" type="text" placeholder="API Key" v-model="apiKey">
        </div>
      </div>
      <div class="field is-grouped">
        <div class="control">
          <button :class="['button', 'is-link']" @click="saveToken">Save</button>
        </div>
        <div class="control">
          <button class="button is-link is-light" @click="toggleApiKeyModal">Cancel</button>
        </div>
      </div>
    </modal>
  </nav>
</template>

<script>
import logo from "../assets/images/logo.png";
import { saveAPIToken } from '../api/auth'
import Modal from '../components/Modal.vue'

export default {
  components: {
    'modal': Modal,
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
    }
  },
  computed: {
    logo_url: function () {
      return logo;
    },
  },
};
</script>
