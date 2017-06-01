require 'capybara'
require 'capybara/cucumber'
require 'capybara/rspec'
require 'cucumber'
require 'fileutils'
require 'git'
require 'pry'
require 'pry-byebug'
require 'selenium-webdriver'

REPO_PATH = '../tmp/repositories/cucumber'
REPO_TEMPLATE_PATH = '../backend/repositories/small'

Capybara.register_driver(:headless_chrome) do |app|
  Capybara::Selenium::Driver.new(
    app,
    browser: :chrome,
    desired_capabilities: Selenium::WebDriver::Remote::Capabilities.chrome(
      chromeOptions: {
        binary: "/Applications/Google\ Chrome\ Canary.app/Contents/MacOS/Google\ Chrome\ Canary",
        args: %w{--headless --no-sandbox}
      }
    )
  )
end

Capybara.register_driver(:chrome) do |app|
  Capybara::Selenium::Driver.new(
    app,
    browser: :chrome,
    desired_capabilities: Selenium::WebDriver::Remote::Capabilities.chrome(
      chromeOptions: {
        binary: "/Applications/Google\ Chrome\ Canary.app/Contents/MacOS/Google\ Chrome\ Canary",
        args: %w{--no-sandbox}
      }
    )
  )
end

Capybara.register_driver :firefox do |app|
	Capybara::Selenium::Driver.new(
		app,
		browser: :firefox,
		desired_capabilities: Selenium::WebDriver::Remote::Capabilities.firefox(marionette: false)
	)
end

Capybara.configure do |c|
  #c.default_driver = :headless_chrome
  #c.default_driver = :firefox
  c.default_driver = :chrome
  c.app_host = "http://localhost:9095"
end

Before do |scenario|

  FileUtils.rm_rf(REPO_PATH)
  FileUtils.cp_r(REPO_TEMPLATE_PATH, REPO_PATH)

  Git.init(REPO_PATH).tap do |g|
    g.add(all: true)
    g.commit("Initial commit")
  end

  @pid = fork do
    %x{../../graphia-cms -config ../../config/cucumber.yml -log-to-file true}
  end

end

After do
  Process.kill 1, @pid
  Process.wait @pid
end

World(Capybara)
