<template>
  <div id="login-form">
    <div class="box">
      <div class="field">
        <label class="label">Email</label>
        <div class="control">
          <input class="input" type="email" placeholder="e.g. alex@example.com" v-model="form.email">
        </div>
      </div>

      <div class="field">
        <label class="label">Password</label>
        <div class="control">
          <input class="input" type="password" placeholder="********" v-model="form.password">
        </div>
      </div>

      <button class="button is-primary" @click="submit">Sign in</button>
    </div>
  </div>
</template>

<script>
import { mapActions } from 'pinia';
import { useUserStore } from '../store/user';
export default {
  data: function () {
    return {
      form: {
        email: "",
        password: "",
      }
    }
  },
  methods: {
    ...mapActions(useUserStore, ['login']),
    submit: async function () {
      let success = await this.login(this.form.email, this.form.password)

      if (success) {
        this.$router.push("/")
      }
    }
  }
};
</script>
