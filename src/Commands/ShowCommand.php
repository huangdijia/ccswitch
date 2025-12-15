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
use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Input\InputOption;
use Symfony\Component\Console\Output\OutputInterface;

#[AsCommand(
    name: 'show',
    description: 'Show details of a specific profile or current settings'
)]
class ShowCommand extends Command
{
    protected function configure(): void
    {
        $this
            ->setHelp('This command allows you to view details of a specific profile or current Claude settings')
            ->addArgument(
                'profile',
                InputArgument::OPTIONAL,
                'The name of the profile to show (omit to show current settings)'
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
            )
            ->addOption(
                'current',
                'c',
                InputOption::VALUE_NONE,
                'Show current Claude settings instead of a profile'
            );
    }

    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $profileName = $input->getArgument('profile');
        $profilesPath = $input->getOption('profiles');
        $showCurrent = $input->getOption('current');

        try {
            // Show current settings if requested or if no profile name provided
            if ($showCurrent || ! $profileName) {
                $profiles = new Profiles($profilesPath);
                $settingsPath = $input->getOption('settings') ?? $profiles->getSettingsPath() ?? '~/.claude/settings.json';
                $currentSettings = new ClaudeSettings($settingsPath);

                $output->writeln('<info>Current Claude Settings:</info>');
                $output->writeln("  Settings file: {$settingsPath}");
                $output->writeln('');

                if (isset($currentSettings->model)) {
                    $output->writeln("  Model: <comment>{$currentSettings->model}</comment>");
                } else {
                    $output->writeln('  Model: <comment>(default)</comment>');
                }

                if (isset($currentSettings->env) && ! empty((array) $currentSettings->env)) {
                    $output->writeln("\n<comment>Environment Variables:</comment>");
                    foreach ((array) $currentSettings->env as $key => $value) {
                        // Hide sensitive information
                        if (stripos($key, 'TOKEN') !== false || stripos($key, 'KEY') !== false) {
                            $value = $this->maskSensitiveValue((string) $value);
                        }
                        $output->writeln("  {$key}: {$value}");
                    }
                }

                return Command::SUCCESS;
            }

            // Show specific profile
            $profiles = new Profiles($profilesPath);

            if (! $profiles->has($profileName)) {
                $output->writeln("<error>Error: Profile '{$profileName}' not found.</error>");
                $output->writeln('<info>Available profiles:</info>');

                foreach (array_keys($profiles->data['profiles'] ?? []) as $name) {
                    $marker = $name === $profiles->default() ? ' *' : '  ';
                    $output->writeln("{$marker}{$name}");
                }

                return Command::FAILURE;
            }

            $profileData = $profiles->get($profileName);
            $descriptions = $profiles->data['descriptions'] ?? [];

            $output->writeln("<info>Profile: {$profileName}</info>");

            if ($profileName === $profiles->default()) {
                $output->writeln('  <comment>(default profile)</comment>');
            }

            if (isset($descriptions[$profileName])) {
                $output->writeln("  Description: {$descriptions[$profileName]}");
            }

            $output->writeln("\n<comment>Configuration:</comment>");

            if (! empty($profileData)) {
                foreach ($profileData as $key => $value) {
                    // Hide sensitive information
                    if (stripos($key, 'TOKEN') !== false || stripos($key, 'KEY') !== false) {
                        $value = $this->maskSensitiveValue((string) $value);
                    }
                    $output->writeln("  {$key}: {$value}");
                }
            } else {
                $output->writeln('  (no custom configuration)');
            }

            return Command::SUCCESS;
        } catch (Exception $e) {
            $output->writeln('<error>Error: ' . $e->getMessage() . '</error>');
            return Command::FAILURE;
        }
    }

    private function maskSensitiveValue(string $value): string
    {
        if (empty($value) || $value === 'sk-' || $value === 'ms-' || $value === 'sk-kimi-') {
            return '(not set)';
        }

        if (strlen($value) <= 8) {
            return str_repeat('*', strlen($value));
        }

        return substr($value, 0, 4) . str_repeat('*', strlen($value) - 8) . substr($value, -4);
    }
}
