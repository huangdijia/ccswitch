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
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Input\InputOption;
use Symfony\Component\Console\Output\OutputInterface;

#[AsCommand(
    name: 'reset',
    description: 'Reset Claude settings to default state'
)]
class ResetCommand extends Command
{
    protected function configure(): void
    {
        $this
            ->setHelp('This command resets your Claude settings to their default state')
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
        try {
            $settingsPath = $input->getOption('settings');

            // If no settings path provided, try to get it from profiles
            if ($settingsPath === null) {
                $profilesPath = getenv('HOME') . '/.ccswitch/ccs.json';
                if (file_exists($profilesPath)) {
                    $profiles = new Profiles($profilesPath);
                    $settingsPath = $profiles->getSettingsPath() ?? '~/.claude/settings.json';
                } else {
                    $settingsPath = '~/.claude/settings.json';
                }
            }

            $settings = new ClaudeSettings($settingsPath);

            // Reset settings to empty state
            $settings->env = new stdClass();

            if (isset($settings->model)) {
                unset($settings->model);
            }

            // Write the reset settings
            $settings->write();

            $output->writeln('<info>Settings have been reset to default</info>');

            return Command::SUCCESS;
        } catch (Exception $e) {
            $output->writeln('<error>Error: ' . $e->getMessage() . '</error>');
            return Command::FAILURE;
        }
    }
}
