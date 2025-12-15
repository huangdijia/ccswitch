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

use Exception;
use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Input\InputOption;
use Symfony\Component\Console\Output\OutputInterface;

#[AsCommand(
    name: 'init',
    description: 'Initialize ccswitch configuration'
)]
class InitCommand extends Command
{
    protected function configure(): void
    {
        $this
            ->setHelp('This command initializes the ccswitch configuration by creating the necessary configuration files and directories')
            ->addOption(
                'profiles',
                'p',
                InputOption::VALUE_OPTIONAL,
                'Path to the profiles configuration file',
                getenv('HOME') . '/.ccswitch/ccs.json'
            )
            ->addOption(
                'force',
                'f',
                InputOption::VALUE_NONE,
                'Force overwrite existing configuration'
            );
    }

    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $profilesPath = $input->getOption('profiles');
        $force = $input->getOption('force');
        $configDir = dirname($profilesPath);

        // Check if configuration already exists
        if (file_exists($profilesPath) && ! $force) {
            $output->writeln('<error>Configuration file already exists. Use --force to overwrite.</error>');
            return Command::FAILURE;
        }

        try {
            // Create config directory if it doesn't exist
            if (! is_dir($configDir)) {
                if (! mkdir($configDir, 0755, true)) {
                    throw new Exception("Failed to create directory: {$configDir}");
                }
                $output->writeln("<info>Created directory: {$configDir}</info>");
            }

            // Default configuration template
            $defaultConfig = [
                'default' => 'default',
                'settingsPath' => getenv('HOME') . '/.config/claude/settings.json',
                'profiles' => [
                    'default' => [
                        'ANTHROPIC_API_KEY' => '',
                        'ANTHROPIC_BASE_URL' => 'https://api.anthropic.com',
                        'ANTHROPIC_MODEL' => 'claude-3-5-sonnet-20241022',
                        'ANTHROPIC_DEFAULT_HAIKU_MODEL' => 'claude-3-5-haiku-20241022',
                        'ANTHROPIC_DEFAULT_OPUS_MODEL' => 'claude-3-opus-20240229',
                        'ANTHROPIC_DEFAULT_SONNET_MODEL' => 'claude-3-5-sonnet-20241022',
                        'ANTHROPIC_SMALL_FAST_MODEL' => 'claude-3-5-haiku-20241022',
                    ],
                ],
                'descriptions' => [
                    'default' => 'Default Claude API configuration',
                ],
            ];

            // Write configuration file
            $jsonContent = json_encode($defaultConfig, JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES);
            if (file_put_contents($profilesPath, $jsonContent) === false) {
                throw new Exception("Failed to write configuration file: {$profilesPath}");
            }

            $output->writeln("<info>Configuration file created successfully: {$profilesPath}</info>");

            return Command::SUCCESS;
        } catch (Exception $e) {
            $output->writeln('<error>Error: ' . $e->getMessage() . '</error>');
            return Command::FAILURE;
        }
    }
}
