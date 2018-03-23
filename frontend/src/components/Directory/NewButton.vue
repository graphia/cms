<template>
	<div class="card">

		<div class="card-body">

			<p>
				A directory is a place to store related documents. You can also add
				more information which will be displayed above the list of contents.
			</p>

			<router-link class="btn btn-lg btn-primary btn-block" :to="{name: 'directory_new'}">
				<octicon :icon-name="'file-directory'"></octicon>
				Create a new directory
			</router-link>
		</div>

	</div>
</template>

<script lang="babel">
	export default {
		name: "DirectoryNewButton"
	}
</script>