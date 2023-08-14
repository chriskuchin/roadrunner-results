import { defineStore, mapActions } from "pinia";
import { useErrorBus } from "../error";
import { setPersistence, browserLocalPersistence, browserSessionPersistence, signInWithEmailAndPassword, createUserWithEmailAndPassword, sendEmailVerification } from 'firebase/auth'
import { auth } from '../../firebase'

export const useUserStore = defineStore('user', {
  state: () => ({
    email: "",
    uid: "",
  }),
  getters: {
    isLoggedIn: (state) => state.uid != "",
    userDisplayName: function(state) {
      if (this.isLoggedIn) {
        return state.email
      }

      return ""
    },
  },
  actions: {
    async register(email, password) {
      try {
        const { user } = await createUserWithEmailAndPassword(auth, email, password)
        await sendEmailVerification(user)
        // create user in main app back end api

      } catch (error) {
        this.handle("Failed to register user: error.code. " + error.message)
      }
    },
    async login(username, password) {
      try {
        await setPersistence(auth, browserLocalPersistence)
        const { user } = await signInWithEmailAndPassword(auth, username, password)

        const token = await user.getIdTokenResult(true)

        console.log(token)
      } catch (error) {
        this.handle("Failed to login. Please verify Username and Password and try again.")
      }

    },
    async logout() {
      await auth.signOut()
      this.email = ""
      this.uid = ""
    },
    async loadUser(fbUser) {
      this.uid = fbUser.uid
      this.email = fbUser.email
    },
    ...mapActions(useErrorBus, ['handle'])
  }
});