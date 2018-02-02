<template>
	<div v-title="this.error.title">
		<div class="alert alert-danger" role="alert">
			<h1>{{ this.code }}</h1>

			<p>
				{{ this.error.message }}
			</p>

			<!--  perhaps add a list of useful links here? -->
		</div>
	</div>
</template>

<script>
	export default {
		name: "Error",
		props: ["code"],
		computed: {
			error() {

				if (!this.code) {
					console.warn("Error component loaded without code");
				};

				return this.errors[this.code];
			}
		},
		data() {
			return {
				errors: {
					404: {
						message: "The item you were looking for cannot be found. Please "      +
								 "ensure you are using the correct URL and that the resource " +
								 "hasn't been deleted.",
						title:   "Not found"
					}
				}
			}
		}
	};
</script>