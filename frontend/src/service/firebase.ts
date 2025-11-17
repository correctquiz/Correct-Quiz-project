// Import the functions you need from the SDKs you need
import { initializeApp } from "firebase/app";
import { getAuth } from "firebase/auth";
// TODO: Add SDKs for Firebase products that you want to use
// https://firebase.google.com/docs/web/setup#available-libraries

// Your web app's Firebase configuration
// For Firebase JS SDK v7.20.0 and later, measurementId is optional
const firebaseConfig = {
  apiKey: "AIzaSyCfzihCuPRE3wwCe96a0I97iroXqtt26Jk",
  authDomain: "correct-quiz-app.firebaseapp.com",
  projectId: "correct-quiz-app",
  storageBucket: "correct-quiz-app.firebasestorage.app",
  messagingSenderId: "368130932433",
  appId: "1:368130932433:web:474a4e7c3e6ced0c304ff0",
  measurementId: "G-FGC86EDJLD"
};

// Initialize Firebase
const app = initializeApp(firebaseConfig);
export const auth = getAuth(app);