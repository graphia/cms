<template>
	<div id="repo-history">

		<Breadcrumbs :levels="breadcrumbs"/>

		<div class="bg-white p-4 m-2">
			<h1>History</h1>

			<div class="commit-list">
				<div class="card m-4" v-for="(commit, i) in commits" :key="i" :data-commit-hash="commit.id">
					<div class="card-header">
						<div class="hash">
							{{ prettyDate(commit.timestamp) }} <span class="text-muted">({{ timeAgo(commit.timestamp) }})</span>
						</div>
						<div class="author">
							{{ commit.author.Name }}
						</div>
					</div>
					<div class="card-body">

						<div v-for="(line, j) in formatMessage(commit.message)" :key="j">
							<h4 v-if="j == 0">
								{{ line }}
							</h4>
							<p v-else>
								{{ line }}
							</p>
						</div>

					</div>

					<div class="card-footer">
						<router-link class="btn btn-sm btn-primary" :to="{name: 'commit', params: {hash: commit.id}}">
							View
						</router-link>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>


<script lang="babel">

	// external
	import fecha from 'fecha';
	import vagueTime from 'vague-time';

	// javascripts
	import config from '../javascripts/config.js';
	import checkResponse from '../javascripts/response.js';
	import CMSBreadcrumb from '../javascripts/models/breadcrumb.js';

	// components
	import Breadcrumbs from './Utilities/Breadcrumbs';


	export default {
		name: "History",
		created() {
			this.getHistory();
		},
		data() {
			return {
				commits: []
			};
		},
		methods: {
			async getHistory() {
				const path = `${config.api}/history`;
				let response = await fetch(path, {headers: this.$store.state.auth.authHeader()})

				if (!checkResponse(response.status)) {
						console.error("History cannot be retrieved", response);
						return;
				};

				this.commits = await response.json();
			},
			parseDate(d) {
				return fecha.parse(d, 'YYYY-MM-DDTHH:mm');
			},
			prettyDate(d) {
				let value = this.parseDate(d);

				return fecha.format(value, 'YYYY-MM-DD HH:mm');
			},
			timeAgo(d) {
				let value = this.parseDate(d);

				return vagueTime.get({
					from: Date.now(),
					to: Date.parse(value),
					units: 'ms'
				});

			},
			formatMessage(m) {
				return m.split(/(\r\n|\n|\r)/gm);
			}
		},
		computed: {
			breadcrumbs() {
				return [new CMSBreadcrumb("History", "history")];
			},
		},
		components: {
			Breadcrumbs
		}
	};
</script>

<style lang="scss" scoped>

	.commit-list {

		.card > .card-header {
			display: flex;

			.hash {
				flex-grow: 1;
			};

		}
	}
</style>