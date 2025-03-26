/// <reference types="svelte" />
/// <reference types="vite/client" />

import { mount } from 'svelte'
import Reset from "./pages/Reset.svelte"
import Auth from "./pages/Auth.svelte"
import Report from "./pages/Report.svelte";
import Settings from "./pages/Settings.svelte";

const pages = new Map()

pages.set("reset", Reset)
pages.set("auth", Auth)
pages.set("report", Report)
pages.set("settings", Settings)

const app = (page: string) => {
    let p = pages.get("*");

    if (pages.has(page)) p = pages.get(page)

    mount(p, {
        target: document.getElementById('app')!,
    })
}

export default app

