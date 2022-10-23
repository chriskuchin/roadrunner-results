<template>
  <div>
    <input
      class="input is-large"
      type="number"
      placeholder="Bib Number Input"
      v-on:keyup.enter="recordRunner"
    />
  </div>
</template>

<script>
export default {
  data: function () {
    return {};
  },
  methods: {
    recordRunner: function (e) {
      e.preventDefault();
      var bibNumber = e.currentTarget.value;
      e.currentTarget.value = "";

      console.log("recordRunner", bibNumber);
    },
  },
};
</script>