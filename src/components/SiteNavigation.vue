<template>
  <b-navbar toggleable="lg" type="dark" variant="info">
    <b-navbar-brand to="/">
      <img
        v-bind:src="logo_url"
        height="30px"
        width="100px"
        alt="RoadRunner Results"
        title="results.roadrunners.club"
      />
    </b-navbar-brand>
    <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>
    <b-collapse id="nav-collapse" is-nav>
      <b-navbar-nav>
        <b-nav-item to="/record">Record</b-nav-item>
        <b-nav-item to="/timer">Timer</b-nav-item>
      </b-navbar-nav>
    </b-collapse>
  </b-navbar>
</template>

<script>
import logo from "../assets/images/logo.png";
import {
  BNavbar,
  BNavbarNav,
  BNavbarBrand,
  BNavbarToggle,
  BNavItem,
  BNavItemDropdown,
  BCollapse,
  BButton,
  BDropdownItem,
  BIcon,
  BIconBug,
  BIconBugFill,
  BIconBag,
  BIconBagFill,
} from "bootstrap-vue";

export default {
  components: {
    "b-navbar": BNavbar,
    "b-navbar-brand": BNavbarBrand,
    "b-navbar-toggle": BNavbarToggle,
    "b-navbar-nav": BNavbarNav,
    "b-nav-item": BNavItem,
    "b-nav-item-dropdown": BNavItemDropdown,
    "b-collapse": BCollapse,
    "b-button": BButton,
    "b-dropdown-item": BDropdownItem,
    "b-icon": BIcon,
    BIconBug,
    BIconBugFill,
    BIconBag,
    BIconBagFill,
  },
  methods: {},
  computed: {
    logo_url: function () {
      return logo;
    },
  },
};
</script>
