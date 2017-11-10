require 'capybara'
require 'capybara/cucumber'
require 'capybara/rspec'
require 'cucumber'
require 'fileutils'
require 'git'
require 'net/http'
require 'pry'
require 'pry-byebug'
require 'selenium-webdriver'

REPO_PATH = '../tmp/repositories/cucumber'
REPO_TEMPLATE_PATH = '../backend/repositories/multiple_filetypes'
PID_PATH = '../tmp/cucumber-browser.pid'
DB_PATH = '../../db/cucumber.db'
SAMPLES_PATH = '../backend/samples'
WEB_SERVER_START_ATTEMPTS = 10 # number of 0.1 second delays to wait for server
DOWNLOAD_DIR = "/tmp/downloads" # needs to be in /tmp because permissions 🚨

DRIVER_PREFS = {
  download: {
    prompt_for_download: false,
    directory_upgrade: true,
    default_directory: DOWNLOAD_DIR
  }
}

Capybara.register_driver(:headless_chrome) do |app|

  Capybara::Selenium::Driver.new(
    app,
    browser: :chrome,
    desired_capabilities: Selenium::WebDriver::Remote::Capabilities.chrome(
      chromeOptions: {
        binary: "/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome",
        args: %w{--headless --no-sandbox --disable-gpu},
        prefs: DRIVER_PREFS
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
        binary: "/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome",
        args: %w{--no-sandbox --disable-gpu},
        prefs: DRIVER_PREFS
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
  if ENV['DRIVER']
    c.default_driver = ENV['DRIVER'].to_sym
  else
    c.default_driver = :headless_chrome
  end

  #c.default_driver = :firefox
  #c.default_driver = :chrome
  c.app_host = "http://localhost:9095"
end

Before do

  # ensure no repo exists
  FileUtils.rm_rf(REPO_PATH)

  # web server start moved to server_steps.rb for extra
  # flexibility

end

After do
  kill(@pid) if @pid
  if FileTest.exist?(DB_PATH)
    File.delete(DB_PATH)
  end
end

def kill(pid)
  Process.kill("HUP", pid)
  Process.wait
  File.delete(PID_PATH)
  @pid = nil
end

World(Capybara)
