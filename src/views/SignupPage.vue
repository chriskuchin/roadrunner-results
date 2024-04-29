<template>
    <div class="section" id="signup-form">
        <div class="box">
            <div class="field">
                <label class="label">Email</label>
                <div class="control">
                    <input class="input" type="email" placeholder="e.g. alex@example.com" v-model="user.username">
                </div>
            </div>

            <div class="field">
                <label class="label">Password</label>
                <div class="control">
                    <input class="input" type="password" placeholder="********" v-model="user.password">
                </div>
            </div>

            <button class="button is-primary" @click="submit">Sign Up</button>
        </div>
    </div>
</template>

<script>
import { mapActions } from 'pinia';
import { useUserStore } from '../store/user';

export default {
    data: function () {
        return {
            user: {
                username: "",
                password: "",
            }
        }
    },
    methods: {
        ...mapActions(useUserStore, ['register']),
        submit: async function () {
            this.register(this.user.username, this.user.password).then(() => {
                this.$router.push('/')
            })
        }
    }
}
</script>