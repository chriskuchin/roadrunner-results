// Import the functions you need from the SDKs you need
import { initializeApp } from "firebase/app";
import { getAuth } from 'firebase/auth'
// TODO: Add SDKs for Firebase products that you want to use
// https://firebase.google.com/docs/web/setup#available-libraries

// Your web app's Firebase configuration
const firebaseConfig = {
  apiKey: "AIzaSyAXgirQia_YKKuK_tHOK1ZLgK1W9qKWp_c",
  authDomain: "rslts-run.firebaseapp.com",
  projectId: "rslts-run",
  storageBucket: "rslts-run.appspot.com",
  messagingSenderId: "601913404882",
  appId: "1:601913404882:web:231b67ef157fac9888ac0b"
};

// Initialize Firebase
const app = initializeApp(firebaseConfig);

const auth = getAuth(app)

export {
  auth
}