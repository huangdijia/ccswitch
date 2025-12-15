<?php

declare(strict_types=1);
/**
 * This file is part of huangdijia/ccswitch.
 *
 * @link     https://github.com/huangdijia/ccswitch
 * @document https://github.com/huangdijia/ccswitch/blob/main/README.md
 * @contact  Your name <your-mail@gmail.com>
 */

namespace CCSwitch\Commands;

use CCSwitch\ClaudeSettings;
use CCSwitch\Profiles;
use Exception;
use stdClass;
use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Input\InputOption;
use Symfony\Component\Console\Output\OutputInterface;

#[AsCommand(
    name: 'set',
    description: 'Set the active Claude API profile'
)]
class SetCommand extends Command
{
    protected function configure(): void
    {
        $this
            ->setHelp('This command allows you to set the active Claude API profile')
            ->addArgument(
                'profile',
                InputArgument::REQUIRED,
                'The name of the profile to set as active'
            )
            ->addOption(
                'profiles',
                'p',
                InputOption::VALUE_OPTIONAL,
                'Path to the profiles configuration file',
                getenv('HOME') . '/.ccswitch/ccs.json'
            )
            ->addOption(
                'settings',
                's',
                InputOption::VALUE_OPTIONAL,
                'Path to the Claude settings file',
                null
            );
    }

    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $profileName = $input->getArgument('profile');
        $profilesPath = $input->getOption('profiles');

        try {
            $profiles = new Profiles($profilesPath);

            if (! $profiles->has($profileName)) {
                $output->writeln("<error>Error: Profile '{$profileName}' not found.</error>");
                $output->writeln('<info>Available profiles:</info>');

                foreach (array_keys($profiles->data['profiles'] ?? []) as $name) {
                    $output->writeln("  - {$name}");
                }

                return Command::FAILURE;
            }

            $settingsPath = $input->getOption('settings') ?? $profiles->getSettingsPath() ?? '~/.claude/settings.json';
            $currentSettings = new ClaudeSettings($settingsPath);

            // Get the environment variables for the selected profile
            $env = $profiles->get($profileName);
            $currentSettings->env = $env ?: new stdClass();

            // Handle model setting
            if (! isset($env['ANTHROPIC_MODEL'])) {
                if (isset($currentSettings->model)) {
                    unset($currentSettings->model);
                }
            } else {
                $currentSettings->model = $env['ANTHROPIC_MODEL'] ?? null;
            }

            // Write settings
            $currentSettings->write();

            $output->writeln("<info>Successfully switched to profile: {$profileName}</info>");

            // Show profile details
            if (! empty($env)) {
                $output->writeln("\n<comment>Profile details:</comment>");
                if (isset($env['ANTHROPIC_BASE_URL'])) {
                    $output->writeln("  URL: {$env['ANTHROPIC_BASE_URL']}");
                }
                if (isset($env['ANTHROPIC_MODEL'])) {
                    $output->writeln("  Model: {$env['ANTHROPIC_MODEL']}");
                }
                if (isset($env['ANTHROPIC_SMALL_FAST_MODEL'])) {
                    $output->writeln("  Fast Model: {$env['ANTHROPIC_SMALL_FAST_MODEL']}");
                }
            }

            return Command::SUCCESS;
        } catch (Exception $e) {
            $output->writeln('<error>Error: ' . $e->getMessage() . '</error>');
            return Command::FAILURE;
        }
    }
}
