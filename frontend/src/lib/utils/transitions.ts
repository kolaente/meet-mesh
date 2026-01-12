import { cubicOut, backOut } from 'svelte/easing'

/**
 * Standard fly in from right
 */
export const flyInRight = { x: 100, duration: 200, easing: cubicOut }

/**
 * Standard fly in from left
 */
export const flyInLeft = { x: -100, duration: 200, easing: cubicOut }

/**
 * Standard fly in from bottom
 */
export const flyInUp = { y: 20, duration: 200, easing: cubicOut }

/**
 * Quick fade for overlays
 */
export const fadeQuick = { duration: 150 }

/**
 * Standard fade for page transitions
 */
export const fadeStandard = { duration: 200 }

/**
 * Subtle scale for buttons/cards on press
 */
export const scalePress = { start: 0.98, duration: 100 }

/**
 * Pop in effect for modals/dialogs
 */
export const popIn = { start: 0.95, duration: 200, easing: backOut }

/**
 * Page transition in
 */
export const pageIn = { duration: 150, delay: 75 }

/**
 * Page transition out
 */
export const pageOut = { duration: 75 }
