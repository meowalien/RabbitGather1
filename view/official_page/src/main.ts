import {createApp} from 'vue'
import App from './App.vue'
import './index.css'
import {store,key} from './store'

const app = createApp(App)

app.use(store,key)

app.mount('#app')
