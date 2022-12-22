import { createApp } from "vue";
import { createPinia } from "pinia";
import { vfmPlugin } from 'vue-final-modal'

import App from "./App.vue";
import router from "./router";

import "./assets/main.css";
import '@mdi/font/css/materialdesignicons.css'

// Vuetify
import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

const vuetify = createVuetify({
    icons: {
        defaultSet: 'mdi', // This is already the default value - only for display purposes
    },
    components,
    directives,
})


const app = createApp(App);

app.use(createPinia());
app.use(router);
app.use(vuetify);
app.use(vfmPlugin);

app.mount("#app");

