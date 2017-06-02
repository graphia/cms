require 'capybara'
require 'capybara/cucumber'
require 'capybara/rspec'
require 'cucumber'
require 'fileutils'
require 'git'
require 'open3'
require 'pry'
require 'pry-byebug'
require 'selenium-webdriver'

REPO_PATH = '../tmp/repositories/cucumber'
REPO_TEMPLATE_PATH = '../backend/repositories/small'
PID_PATH = '../tmp/cucumber-browser.pid'

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
  c.default_driver = :headless_chrome
  #c.default_driver = :firefox
  #c.default_driver = :chrome
  c.app_host = "http://localhost:9095"
end

Before do |scenario|

  FileUtils.rm_rf(REPO_PATH)
  FileUtils.cp_r(REPO_TEMPLATE_PATH, REPO_PATH)

  Git.init(REPO_PATH).tap do |g|
    g.add(all: true)
    g.commit("Initial commit")
  end

  # kill existing pid first
  if FileTest.exist?(PID_PATH)
    begin
      @pid = Pathname.new(PID_PATH).read.to_i
      kill(@pid)
    rescue Errno::ESRCH
      # already dead 😵
    end
  end

  @pid = fork do
    #system("../../graphia-cms -config ../../config/cucumber.yml -log-to-file true")
    command = "../../graphia-cms -config=../../config/cucumber.yml -log-to-file=true "

    Open3.popen3(command) do |stdin, stdout, stderr, wait_thr|
      # FIXME Negroni's output is still appearing, work out how to suppress it
    end

  end

  Pathname.new(PID_PATH).write(@pid)

end

After do
  kill(@pid)
end

def kill(pid)
  Process.kill 9, pid
  Process.wait pid
  File.delete(PID_PATH)
  @pid = nil
end

World(Capybara)