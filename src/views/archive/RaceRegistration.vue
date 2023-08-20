<template>
  <div class="section">
    <h1 class="title">Registration</h1>

  </div>
</template>

<script>
export default {
  data: function () {
    return {
      participant: {
        bib_number: "",
        first_name: "",
        last_name: "",
        gender: "",
        birth_year: new Date().getFullYear(),
        team: "",
      }
    }
  },
  method: {
    createParticipant: function () {

    },
    clearForm: function () {

    }
  },
  computed: {
  }
}
</script>