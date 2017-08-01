<template>

	<div class="commit-message form-group">
		<label for="commit-message">Commit Message</label>
		<textarea name="commit-message" class="form-control" v-model="commit.message"/>
	</div>

</template>

<script lang="babel">
	export default {
		name: "CommitMessage",
		computed: {
			commit() {
				return this.$store.state.commit;
			}
		}
	}
</script>