import firebase from 'firebase'

export const db = firebase
  .initializeApp({ projectId: 'bostadskoe' })
  .firestore()

export const { TimeStamp, GeoPoint } = firebase.firestore
