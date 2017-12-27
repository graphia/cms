<template>
	<div id="broadcast-messages" class="row messages">
		<div class="col-md-12 mt-2">
			<div v-for="(message, i) in messages" :key=i :class="message.classes" role="alert">
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
				if (this.$store.state.broadcast) {
					return this.$store.state.broadcast.activeMessages();
				};

				return [];
			}
		}
	};
</script>