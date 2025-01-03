# Wire Dependency Injection Notes

## Key Issues Encountered

1. **Wire Injector Function Restrictions**
   - Wire injector functions (those using `wire.Build`) must ONLY contain the wire.Build call
   - No debug prints or other code allowed in these functions
   - Example of what NOT to do:
   ```go
   func InitializeTestApp() (*app.TestApp, func(), error) {
       fmt.Println("Debug...") // NOT ALLOWED
       wire.Build(TestSet)
       return nil, nil, nil
   }
   ```

2. **Provider Set Organization**
   - Keep provider sets focused and small
   - Current structure:
     - `CoreSet`: Essential services (logging, storage, config)
     - `UISet`: UI components
     - `HotkeySet`: Platform features
     - `HTTPSet`: Server config
     - `AppSet`: Main app wiring
     - `TestSet`: Mock dependencies

3. **Mock Provider Issues**
   - Mock providers need careful binding to interfaces
   - Current bindings in TestSet:
   ```go
   wire.Bind(new(gui.MainWindow), new(*gui.MockMainWindow))
   wire.Bind(new(gui.QuickNote), new(*gui.MockQuickNote))
   wire.Bind(new(apphotkey.Manager), new(*apphotkey.MockManager))
   ```

4. **Debugging Strategies**
   - Can't add debug prints to injector functions
   - Instead, add logging to provider functions
   - Use separate debug provider functions for testing
   - Example:
   ```go
   func ProvideDebugMockStore() storage.Store {
       result := ProvideMockStore()
       fmt.Printf("Debug: Providing MockStore: %T %+v\n", result, result)
       return result
   }
   ```

## Current Test Issues

1. `TestNew` in `container_test.go` fails with nil TestApp
   - Possible causes:
     - Missing provider in TestSet
     - Incorrect interface binding
     - Provider returning nil
   - Need to verify all mock providers are properly initialized

## Next Steps

1. Consider removing or simplifying the container tests
2. If keeping tests:
   - Move debug logging to provider functions
   - Add validation in mock providers
   - Consider simpler test structure

## Wire Best Practices

1. Keep injector functions clean (wire.Build only)
2. Use provider sets for organization
3. Validate dependencies in provider functions
4. Use proper interface bindings
5. Test with focused mock providers 