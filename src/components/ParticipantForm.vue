<template>
  <div class="container">
    <div class="field is-horizontal mb-2 mt-4">
      <div class="field-label is-normal">
        <label class="label">Name</label>
      </div>
      <div class="field-body">
        <div class="field">
          <p class="control is-expanded">
            <input class="input" type="text" placeholder="First Name" :value="first_name"
              @input="$emit('update:first_name', $event.target.value)">
          </p>
        </div>
        <div class="field">
          <p class="control is-expanded">
            <input class="input" type="tex" placeholder="Last Name" :value="last_name"
              @input="$emit('update:last_name', $event.target.value)">
          </p>
        </div>
      </div>
    </div>

    <div class="field is-horizontal mb-2">
      <div class="field-label is-normal">
        <label class="label">Race Info</label>
      </div>
      <div class="field-body">
        <div class="field">
          <p class="control is-expanded">
            <input class="input" type="text" placeholder="Bib Number" :value="bib_number"
              :disabled="disable_bib_input != undefined" @input="$emit('update:bib_number', $event.target.value)">
          </p>
        </div>
        <div class="field">
          <p class="control is-expanded">
            <input class="input" type="text" placeholder="Team" :value="team"
              @input="$emit('update:team', $event.target.value)">
          </p>
        </div>
      </div>
    </div>

    <div class="field is-horizontal mb-2">
      <div class="field-label is-normal">
        <label class="label">Gender</label>
      </div>
      <div class="field-body">
        <div class="field">
          <div class="control is-expanded">
            <label class="radio">
              <input type="radio" name="gender" value="Male" :checked="gender == 'Male'"
                @input="$emit('update:gender', $event.target.value)">
              Male
            </label>
            <label class="radio">
              <input type="radio" name="gender" value="Female" :checked="gender == 'Female'"
                @input="$emit('update:gender', $event.target.value)">
              Female
            </label>
          </div>
        </div>
      </div>
    </div>

    <div class="field is-horizontal mb-2">
      <div class="field-label is-normal">
        <label class="label">Birth Year</label>
      </div>
      <div class="field-body">
        <div class="field">
          <div class="control is-expanded">
            <div class="select">
              <select :value="birth_year" @input="$emit('update:birth_year', $event.target.value)">
                <option>{{ year }}</option>
                <option v-for=" n  in  80 ">{{ year - n }}</option>
              </select>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>

export default {
  props: ['first_name', 'last_name', 'bib_number', 'birth_year', 'team', 'gender', 'disable_bib_input'],
  emits: ['update:first_name', 'update:last_name', 'update:bib_number', 'update:birth_year', 'update:team', 'update:gender'],
  data: function () {
    return {}
  },
  methods: {

  },
  computed: {
    year: () => new Date().getFullYear(),
  }
}
</script>