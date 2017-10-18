<script lang="babel">
	export default {
		name: "DocumentConflict",
		computed: {
			document() {
				return this.$store.state.activeDocument;
			},
			directory() {
				return this.$route.params.directory;
			}
		},
		methods: {
			downloadFile() {

				let a = document.createElement("a");

				document.body.appendChild(a);
				a.style = "display: none";

				let blob = new Blob(
						[this.document.markdown],
						{type: "octet/stream"}
					),
					url = window.URL.createObjectURL(blob);

				a.href = url;
				a.download = this.document.filename;
				a.click();
				window.URL.revokeObjectURL(url);

			}
		}
	};
</script>

<template>

	<div id="conflict-warning" class="modal fade">
		<div class="modal-dialog" role="document">
			<div class="modal-content">
				<div class="modal-header">
					<h5 class="modal-title">Your commit failed</h5>
					<button type="button" class="close" data-dismiss="modal" aria-label="Close">
					<span aria-hidden="true">&times;</span>
					</button>
				</div>
				<div class="modal-body">
					<p>
						The repository has been modified since you began editing.
					</p>
				</div>
				<div class="modal-footer">
					<button @click="downloadFile" type="button" class="btn btn-success">Download your copy</button>
					<router-link class="btn btn-danger" data-dismiss="modal" :to="{name: 'document_index', params: {directory: this.directory}}">
						Close
					</router-link>
				</div>
			</div>
		</div>
	</div>

</template>