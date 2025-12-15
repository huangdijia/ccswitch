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
use Symfony\Component\Console\Style\SymfonyStyle;

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
                'full',
                null,
                InputOption::VALUE_NONE,
                'Use full configuration with all available providers'
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
        $full = $input->getOption('full');
        $configDir = dirname($profilesPath);
        $io = new SymfonyStyle($input, $output);

        // Check if configuration already exists
        if (file_exists($profilesPath) && ! $force) {
            $io->error('Configuration file already exists. Use --force to overwrite.');
            return Command::FAILURE;
        }

        try {
            // Create config directory if it doesn't exist
            if (! is_dir($configDir)) {
                if (! mkdir($configDir, 0755, true)) {
                    throw new Exception("Failed to create directory: {$configDir}");
                }
                $io->success("Created directory: {$configDir}");
            }

            // Determine source config file
            $projectRoot = dirname(__DIR__, 2);
            $sourceConfig = $full
                ? $projectRoot . '/config/ccs-full.json'
                : $projectRoot . '/config/ccs.json';

            // Check if source config exists
            if (! file_exists($sourceConfig)) {
                $io->error("Source configuration file not found: {$sourceConfig}");
                return Command::FAILURE;
            }

            // Copy configuration file
            $configContent = file_get_contents($sourceConfig);
            if ($configContent === false) {
                throw new Exception("Failed to read source configuration: {$sourceConfig}");
            }

            if (file_put_contents($profilesPath, $configContent) === false) {
                throw new Exception("Failed to write configuration file: {$profilesPath}");
            }

            $configType = $full ? 'full' : 'default';
            $io->success("{$configType} configuration file created successfully: {$profilesPath}");

            return Command::SUCCESS;
        } catch (Exception $e) {
            $io->error('Error: ' . $e->getMessage());
            return Command::FAILURE;
        }
    }
}
