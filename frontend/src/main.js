import { createRouter, createWebHistory } from 'vue-router';
import { createApp } from 'vue';

import Dashboard from './Dashboard.vue';
import Connections from './Connections.vue';
import App from './App.vue';

import 'bootstrap/dist/js/bootstrap.bundle.min.js';
import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap-icons/font/bootstrap-icons.min.css';

const app = createApp(App);

const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/',
            name: 'home',
            component: Connections
        },
        {
            path: '/dashboard',
            name: 'dashboard',
            component: Dashboard
        }
    ]
});

app.use(router);
app.mount('#app');
