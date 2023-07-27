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
            <input class="input" type="text" placeholder="First Name" v-model="participant.first_name">
          </p>
        </div>
        <div class="field">
          <p class="control is-expanded">
            <input class="input" type="tex" placeholder="Last Name" v-model="participant.last_name">
          </p>
        </div>
      </div>
    </div>

    <div class="field is-horizontal">
      <div class="field-label is-normal">
        <label class="label">Race Info</label>
      </div>
      <div class="field-body">
        <div class="field">
          <p class="control is-expanded">
            <input class="input" type="text" placeholder="Bib Number" v-model="participant.bib_number">
          </p>
        </div>
        <div class="field">
          <p class="control is-expanded">
            <input class="input" type="text" placeholder="Team" v-model="participant.team">
          </p>
        </div>
      </div>
    </div>

    <div class="field is-horizontal">
      <div class="field-label is-normal">
        <label class="label">Gender</label>
      </div>
      <div class="field-body">
        <div class="field">
          <div class="control is-expanded">
            <label class="radio">
              <input type="radio" name="gender" v-model="participant.gender">
              Male
            </label>
            <label class="radio">
              <input type="radio" name="gender" v-model="participant.gender">
              Female
            </label>
          </div>
        </div>
      </div>
    </div>

    <div class="field is-horizontal">
      <div class="field-label is-normal">
        <label class="label">Birth Year</label>
      </div>
      <div class="field-body">
        <div class="field">
          <div class="control is-expanded">
            <div class="select">
              <select v-model="participant.birth_year">
                <option>{{ year }}</option>
                <option v-for="n in 80">{{ year - n }}</option>
              </select>
            </div>
          </div>
        </div>
      </div>
    </div>


    <div class="field is-horizontal">
      <div class="field-label">
        <!-- Left empty for spacing -->
      </div>
      <div class="field-body">
        <div class="field is-grouped">
          <div class="control">
            <button class="button is-link" @click="createParticipant">Submit</button>
          </div>
          <div class="control">
            <button class="button is-link is-light" @click="clearForm">Cancel</button>
          </div>
        </div>
      </div>
    </div>
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
    year: () => new Date().getFullYear(),
    raceID: () => this.$route.params.raceId,
  }
}
</script>