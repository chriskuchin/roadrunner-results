<template>
  <div class="section">
    <h1 class="title">Registration</h1>

    <div class="field is-horizontal">
      <div class="field-label is-normal">
        <label class="label">Name</label>
      </div>
      <div class="field-body">
        <div class="field">
          <p class="control is-expanded">
            <input class="input" type="text" placeholder="First Name">
          </p>
        </div>
        <div class="field">
          <p class="control is-expanded">
            <input class="input" type="tex" placeholder="Last Name">
          </p>
        </div>
      </div>
    </div>

    <div class="field">
      <label class="label">Bib Number</label>
      <div class="control">
        <input class="input" type="text" placeholder="Text input">
      </div>
    </div>

    <div class="field">
      <label class="label">First Name</label>
      <div class="control">
        <input class="input" type="text" placeholder="Text input">
      </div>
    </div>

    <div class="field">
      <label class="label">Last Name</label>
      <div class="control">
        <input class="input" type="text" placeholder="Text input">
      </div>
    </div>

    <div class="field">
      <label class="label">Team</label>
      <div class="control">
        <input class="input" type="text" placeholder="Text input">
      </div>
    </div>

    <div class="field">
      <label class="label">Birth Year</label>
      <div class="control">
        <div class="select">
          <select>
            <option>{{ year }}</option>
            <option v-for="n in 80">{{ year - n }}</option>
          </select>
        </div>
      </div>
    </div>

    <div class="field">
      <div class="control">
        <label class="radio">
          <input type="radio" name="gender">
          Male
        </label>
        <label class="radio">
          <input type="radio" name="gender">
          Female
        </label>
      </div>
    </div>

    <div class="field is-grouped">
      <div class="control">
        <button class="button is-link">Submit</button>
      </div>
      <div class="control">
        <button class="button is-link is-light">Cancel</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  date: function () {
    return {
      participant: {
        bib_number: "",
        first_name: "",
        last_name: "",
        gender: "",
        birth_year: "",
        team: "",
      }
    }
  },
  method: {
  },
  computed: {
    year: () => new Date().getFullYear(),
    raceID: () => this.$route.params.raceId,
  }
}
</script>