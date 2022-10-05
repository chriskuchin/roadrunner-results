<template>
  <nav class="navbar is-dark" role="navigation" aria-label="main navigation">
    <div class="navbar-brand">
      <a class="navbar-item" href="/">
        <img
          v-bind:src="logo_url"
          height="30px"
          alt="RoadRunner Results"
          title="results.roadrunners.club"
        />
      </a>
      <a
        role="button"
        class="navbar-burger"
        aria-label="menu"
        aria-expanded="false"
        data-target="navbarBasicExample"
        @click="toggleBurger"
        :class="{ 'is-active': activeBurger }"
      >
        <span aria-hidden="true"></span>
        <span aria-hidden="true"></span>
        <span aria-hidden="true"></span>
      </a>
    </div>
    <div
      id="navbarBasicExample"
      class="navbar-menu"
      :class="{ 'is-active': activeBurger }"
    >
      <div class="navbar-start">
        <router-link to="/record" class="navbar-item">Record</router-link>
        <router-link to="/timer" class="navbar-item">Timer</router-link>
      </div>
      <div class="navbar-end">
        <div class="navbar-item">
          <div class="buttons">
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
  </nav>
</template>

<script>
import logo from "../assets/images/logo.png";

export default {
  data: function () {
    return {
      activeBurger: false,
    };
  },
  components: {},
  methods: {
    toggleBurger: function () {
      this.activeBurger = !this.activeBurger;
    },
  },
  computed: {
    logo_url: function () {
      return logo;
    },
  },
};
</script>
