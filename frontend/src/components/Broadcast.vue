<template>
	<div class="messages">
		<div v-for="message in messages" :class="message.classes" role="alert">
			{{ message.content }}
		</div>
	</div>
</template>

<script lang="babel">
	export default {
		name: "Broadcast",
		computed: {
			messages() {
				return this.$store.state.broadcast.activeMessages();
			}
		}
	};
</script>