<template>
	<div class="row messages">
		<div class="col-md-12 mt-2">
			<div v-for="message in messages" :class="message.classes" role="alert">
				<strong>{{ message.alert }}</strong>
				{{ message.content }}
			</div>
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