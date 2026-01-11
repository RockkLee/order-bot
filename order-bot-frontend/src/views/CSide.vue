<script setup lang="ts">
import { ref } from 'vue'

type Panel = 'dialogue' | 'menu'

const activePanel = ref<Panel>('dialogue')

const setPanel = (panel: Panel) => {
  activePanel.value = panel
}
</script>

<template>
  <section class="panel-shell">
    <header class="panel-top">
      <div>
        <p class="eyebrow">C-side</p>
        <h1>Guest experience</h1>
        <p class="panel-subtitle">No login needed. Tap a view to switch.</p>
      </div>
      <div class="panel-actions">
        <button
          type="button"
          :class="['toggle-btn', activePanel === 'dialogue' ? 'is-active' : '']"
          @click="setPanel('dialogue')"
        >
          Dialogue
        </button>
        <button
          type="button"
          :class="['toggle-btn', activePanel === 'menu' ? 'is-active' : '']"
          @click="setPanel('menu')"
        >
          Show the menu
        </button>
      </div>
    </header>

    <div class="panel-content">
      <div v-if="activePanel === 'dialogue'" class="dialogue-view">
        <div class="chat-bubble incoming">
          <p>Welcome! I can help you pick today&apos;s special.</p>
          <span>Order Bot</span>
        </div>
        <div class="chat-bubble outgoing">
          <p>Show me something light and spicy.</p>
          <span>You</span>
        </div>
        <div class="chat-bubble incoming">
          <p>Try the Citrus Chili Bowl with sparkling yuzu.</p>
          <span>Order Bot</span>
        </div>
        <div class="chat-bubble outgoing">
          <p>Add it to my order.</p>
          <span>You</span>
        </div>
        <div class="chat-input">
          <input placeholder="Ask for a recommendation..." />
          <button type="button">Send</button>
        </div>
      </div>

      <div v-else class="menu-view">
        <div class="menu-card">
          <h2>Today&apos;s menu</h2>
          <ul>
            <li>
              <span>Citrus Chili Bowl</span>
              <span>$12</span>
            </li>
            <li>
              <span>Miso Maple Noodles</span>
              <span>$11</span>
            </li>
            <li>
              <span>Charred Corn Salad</span>
              <span>$8</span>
            </li>
            <li>
              <span>Cold Brew Spritz</span>
              <span>$5</span>
            </li>
          </ul>
          <button type="button" class="primary-btn">Start an order</button>
        </div>
        <div class="menu-highlight">
          <p class="eyebrow">Popular right now</p>
          <h3>+120 orders today</h3>
          <p>Personalize spice, swap grains, and save favorites.</p>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.panel-shell {
  display: flex;
  flex-direction: column;
  gap: 24px;
  padding: 28px;
  border-radius: var(--radius-lg);
  background: var(--card);
  box-shadow: var(--shadow);
}

.panel-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
  flex-wrap: wrap;
}

h1 {
  margin: 6px 0 4px;
  font-size: clamp(1.6rem, 2.6vw, 2.4rem);
}

.panel-subtitle {
  margin: 0;
  color: var(--muted);
}

.panel-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.toggle-btn {
  border: none;
  padding: 10px 18px;
  border-radius: 999px;
  background: #f0e6d7;
  color: var(--ink);
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.toggle-btn.is-active {
  background: var(--accent);
  color: #fff7ed;
  transform: translateY(-1px);
  box-shadow: 0 10px 18px rgba(210, 79, 35, 0.35);
}

.panel-content {
  display: grid;
  gap: 20px;
}

.dialogue-view {
  display: grid;
  gap: 16px;
}

.chat-bubble {
  max-width: 420px;
  padding: 14px 16px;
  border-radius: 18px;
  display: grid;
  gap: 6px;
  position: relative;
  box-shadow: 0 12px 20px rgba(19, 32, 28, 0.15);
}

.chat-bubble span {
  font-size: 0.75rem;
  color: var(--muted);
}

.incoming {
  background: #ffffff;
}

.outgoing {
  margin-left: auto;
  background: #f7d8c8;
}

.chat-input {
  display: flex;
  gap: 12px;
  padding: 12px;
  background: #fff5e8;
  border-radius: 16px;
}

.chat-input input {
  flex: 1;
  border: none;
  background: transparent;
  font-size: 1rem;
  outline: none;
}

.chat-input button {
  border: none;
  padding: 10px 16px;
  border-radius: 999px;
  background: var(--accent);
  color: #fff7ed;
  font-weight: 600;
  cursor: pointer;
}

.menu-view {
  display: grid;
  gap: 20px;
  grid-template-columns: minmax(0, 2fr) minmax(0, 1fr);
}

.menu-card {
  background: #fffaf2;
  border-radius: 20px;
  padding: 20px 22px;
  display: grid;
  gap: 16px;
}

.menu-card ul {
  list-style: none;
  padding: 0;
  margin: 0;
  display: grid;
  gap: 12px;
}

.menu-card li {
  display: flex;
  justify-content: space-between;
  font-weight: 600;
}

.primary-btn {
  border: none;
  padding: 12px 18px;
  border-radius: 999px;
  background: #13201c;
  color: #fff7ed;
  font-weight: 600;
  cursor: pointer;
}

.menu-highlight {
  border-radius: 20px;
  padding: 20px;
  background: linear-gradient(135deg, #d5f7e4 0%, #fffaf2 100%);
  display: grid;
  gap: 8px;
}

.eyebrow {
  margin: 0;
  font-size: 0.75rem;
  letter-spacing: 0.2em;
  text-transform: uppercase;
  color: var(--muted);
}

@media (max-width: 900px) {
  .menu-view {
    grid-template-columns: 1fr;
  }
}
</style>
