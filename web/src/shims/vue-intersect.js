// This is a patched version of vue-intersect to work with Vue 3
import { defineComponent, h, ref, onMounted, onBeforeUnmount } from 'vue'

export default defineComponent({
  name: 'Intersect',
  props: {
    threshold: {
      type: Array,
      required: false,
      default: () => [0, 0.2]
    },
    root: {
      type: Object,
      required: false,
      default: null
    },
    rootMargin: {
      type: String,
      required: false,
      default: '0px 0px 0px 0px'
    }
  },
  emits: ['enter', 'leave', 'change'],
  setup(props, { slots, emit }) {
    const observer = ref(null)
    const isIntersecting = ref(false)
    const element = ref(null)

    const createObserver = () => {
      if (!window.IntersectionObserver) return
      
      // Destroy existing observer if it exists
      if (observer.value) {
        destroyObserver()
      }

      // Create new observer
      observer.value = new IntersectionObserver(
        (entries) => {
          const isIntersectingNow = entries[0].isIntersecting
          
          if (isIntersectingNow !== isIntersecting.value) {
            isIntersecting.value = isIntersectingNow
            
            if (isIntersectingNow) {
              emit('enter', entries[0])
            } else {
              emit('leave', entries[0])
            }
            
            emit('change', isIntersectingNow, entries[0])
          }
        },
        {
          threshold: props.threshold,
          root: props.root,
          rootMargin: props.rootMargin
        }
      )

      if (element.value) {
        observer.value.observe(element.value)
      }
    }

    const destroyObserver = () => {
      if (observer.value) {
        observer.value.disconnect()
        observer.value = null
      }
    }

    onMounted(() => {
      createObserver()
    })

    onBeforeUnmount(() => {
      destroyObserver()
    })

    return () => {
      // Render the default slot and capture the element reference
      const content = slots.default ? slots.default({ isIntersecting: isIntersecting.value }) : []
      
      if (content.length === 1 && content[0].type !== Comment) {
        const vnode = content[0]
        const originalRef = vnode.ref
        
        vnode.ref = (el) => {
          element.value = el
          
          // Call the original ref if it exists
          if (originalRef) {
            if (typeof originalRef === 'function') {
              originalRef(el)
            } else if (originalRef.hasOwnProperty('current')) {
              originalRef.current = el
            }
          }
          
          // If we already have an observer, observe this element
          if (observer.value && el) {
            observer.value.observe(el)
          }
        }
        
        return vnode
      } else {
        // Wrap multiple children in a div
        return h('div', { ref: (el) => { element.value = el } }, content)
      }
    }
  }
})
